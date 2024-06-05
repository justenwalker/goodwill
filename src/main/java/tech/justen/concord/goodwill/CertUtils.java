// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import java.io.*;
import java.math.BigInteger;
import java.security.*;
import java.security.cert.CertificateEncodingException;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.security.spec.ECGenParameterSpec;
import java.util.Calendar;
import java.util.Date;
import org.bouncycastle.asn1.ASN1Encodable;
import org.bouncycastle.asn1.DERSequence;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.Extension;
import org.bouncycastle.asn1.x509.GeneralName;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509ExtensionUtils;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.openssl.jcajce.JcaPEMWriter;
import org.bouncycastle.openssl.jcajce.JcaPKCS8Generator;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.bouncycastle.pkcs.PKCS10CertificationRequest;
import org.bouncycastle.pkcs.PKCS10CertificationRequestBuilder;
import org.bouncycastle.pkcs.jcajce.JcaPKCS10CertificationRequestBuilder;
import org.bouncycastle.util.io.pem.PemObject;

public class CertUtils {
  public static class CA {
    private X509CertificateHolder cert;

    private KeyPair keyPair;

    public CA(KeyPair keyPair, X509CertificateHolder cert) {
      this.cert = cert;
      this.keyPair = keyPair;
    }

    public InputStream getCACertInputStream() throws IOException, CertificateException {
      X509Certificate ca =
          new JcaX509CertificateConverter().setProvider(BC_PROVIDER).getCertificate(this.cert);
      ByteArrayOutputStream baos = new ByteArrayOutputStream();
      try (OutputStreamWriter osw = new OutputStreamWriter(baos)) {
        encodeCertificate(ca, osw);
      }
      return new ByteArrayInputStream(baos.toByteArray());
    }

    public InputStream getCAKeyInputStream() throws IOException {
      ByteArrayOutputStream baos = new ByteArrayOutputStream();
      try (OutputStreamWriter osw = new OutputStreamWriter(baos)) {
        encodePrivateKey(keyPair.getPrivate(), osw);
      }
      return new ByteArrayInputStream(baos.toByteArray());
    }

    public void generatePKI(File caCert, File certFile, File keyFile) throws Exception {
      X509Certificate ca =
          new JcaX509CertificateConverter().setProvider(BC_PROVIDER).getCertificate(this.cert);
      encodeCertificate(ca, caCert);

      X500Name issuedCertSubject = new X500Name("CN=goodwill-client");
      BigInteger serial = new BigInteger(Long.toString(new SecureRandom().nextLong()));
      KeyPair issuedCertKeyPair = generateKeyPair();

      encodePrivateKey(issuedCertKeyPair.getPrivate(), keyFile);

      Calendar calendar = Calendar.getInstance();
      Date startDate = calendar.getTime();

      calendar.add(Calendar.DATE, 1);
      Date endDate = calendar.getTime();

      PKCS10CertificationRequestBuilder p10Builder =
          new JcaPKCS10CertificationRequestBuilder(
              issuedCertSubject, issuedCertKeyPair.getPublic());
      JcaContentSignerBuilder csrBuilder =
          new JcaContentSignerBuilder(SIGNATURE_ALGORITHM).setProvider(BC_PROVIDER);

      // Sign the new KeyPair with the root cert Private Key
      ContentSigner csrContentSigner = csrBuilder.build(keyPair.getPrivate());
      PKCS10CertificationRequest csr = p10Builder.build(csrContentSigner);

      X509v3CertificateBuilder issuedCertBuilder =
          new X509v3CertificateBuilder(
              cert.getSubject(),
              serial,
              startDate,
              endDate,
              csr.getSubject(),
              csr.getSubjectPublicKeyInfo());

      JcaX509ExtensionUtils issuedCertExtUtils = new JcaX509ExtensionUtils();

      // Add Extensions
      // Use BasicConstraints to say that this Cert is not a CA
      issuedCertBuilder.addExtension(Extension.basicConstraints, true, new BasicConstraints(false));

      // Add Issuer cert identifier as Extension
      issuedCertBuilder.addExtension(
          Extension.authorityKeyIdentifier,
          false,
          issuedCertExtUtils.createAuthorityKeyIdentifier(cert));
      issuedCertBuilder.addExtension(
          Extension.subjectKeyIdentifier,
          false,
          issuedCertExtUtils.createSubjectKeyIdentifier(csr.getSubjectPublicKeyInfo()));

      // Add intended key usage extension if needed
      issuedCertBuilder.addExtension(
          Extension.keyUsage, false, new KeyUsage(KeyUsage.keyEncipherment));

      X509CertificateHolder issuedCertHolder = issuedCertBuilder.build(csrContentSigner);
      X509Certificate cert =
          new JcaX509CertificateConverter()
              .setProvider(BC_PROVIDER)
              .getCertificate(issuedCertHolder);
      encodeCertificate(cert, certFile);
    }
  }

  private static final String BC_PROVIDER = "BC";

  private static final String KEY_ALGORITHM = "EC";

  private static final String SIGNATURE_ALGORITHM = "SHA256withECDSA";

  private static final ECGenParameterSpec EC_PARAMS = new ECGenParameterSpec("secp256r1");

  private static void registerProvider() {
    if (Security.getProvider(BC_PROVIDER) == null) {
      Security.addProvider(new BouncyCastleProvider());
    }
  }

  public static KeyPair generateKeyPair()
      throws NoSuchProviderException, NoSuchAlgorithmException, InvalidAlgorithmParameterException {
    registerProvider();
    KeyPairGenerator keygen = KeyPairGenerator.getInstance(KEY_ALGORITHM, BC_PROVIDER);
    keygen.initialize(EC_PARAMS);
    return keygen.generateKeyPair();
  }

  public static void encodeCertificate(X509Certificate cert, File file)
      throws IOException, CertificateEncodingException {
    try (FileWriter fw = new FileWriter(file)) {
      encodeCertificate(cert, fw);
    }
  }

  public static void encodeCertificate(X509Certificate cert, Writer writer)
      throws IOException, CertificateEncodingException {
    try (JcaPEMWriter pw = new JcaPEMWriter(writer)) {
      pw.writeObject(new PemObject("CERTIFICATE", cert.getEncoded()));
    }
  }

  public static void encodePrivateKey(PrivateKey key, File file) throws IOException {
    try (FileWriter fw = new FileWriter(file)) {
      encodePrivateKey(key, fw);
    }
  }

  public static void encodePrivateKey(PrivateKey key, Writer writer) throws IOException {
    JcaPKCS8Generator gen1 = new JcaPKCS8Generator(key, null);
    try (JcaPEMWriter pw = new JcaPEMWriter(writer)) {
      pw.writeObject(gen1.generate());
    }
  }

  public static CA generateCA() throws Exception {
    KeyPair rootKeyPair = generateKeyPair();

    Calendar calendar = Calendar.getInstance();
    Date startDate = calendar.getTime();

    calendar.add(Calendar.DATE, 3);
    Date endDate = calendar.getTime();

    BigInteger rootSerialNum = new BigInteger(Long.toString(new SecureRandom().nextLong()));
    X500Name rootCertIssuer = new X500Name("CN=goodwill-server");
    X500Name rootCertSubject = rootCertIssuer;
    ContentSigner rootCertContentSigner =
        new JcaContentSignerBuilder(SIGNATURE_ALGORITHM)
            .setProvider(BC_PROVIDER)
            .build(rootKeyPair.getPrivate());
    X509v3CertificateBuilder rootCertBuilder =
        new JcaX509v3CertificateBuilder(
            rootCertIssuer,
            rootSerialNum,
            startDate,
            endDate,
            rootCertSubject,
            rootKeyPair.getPublic());

    JcaX509ExtensionUtils rootCertExtUtils = new JcaX509ExtensionUtils();
    rootCertBuilder.addExtension(Extension.basicConstraints, true, new BasicConstraints(true));
    rootCertBuilder.addExtension(
        Extension.subjectKeyIdentifier,
        false,
        rootCertExtUtils.createSubjectKeyIdentifier(rootKeyPair.getPublic()));
    rootCertBuilder.addExtension(
        Extension.subjectAlternativeName,
        false,
        new DERSequence(
            new ASN1Encodable[] {
              new GeneralName(GeneralName.dNSName, "localhost"),
              new GeneralName(GeneralName.iPAddress, "127.0.0.1")
            }));
    X509CertificateHolder rootCertHolder = rootCertBuilder.build(rootCertContentSigner);

    return new CA(rootKeyPair, rootCertHolder);
  }
}
