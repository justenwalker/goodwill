// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

import java.io.File;

public class CertUtilsTest {

    @Rule
    public TemporaryFolder folder = new TemporaryFolder();

    @Test
    public void testGenerateCA() throws Exception {
        File caFile = folder.newFile("ca.crt");
        File certFile = folder.newFile("client.crt");
        File keyFile = folder.newFile("client.key");
        CertUtils.CA ca = CertUtils.generateCA();
        ca.generatePKI(caFile, certFile, keyFile);
        Assert.assertTrue(caFile.exists());
        Assert.assertTrue(certFile.exists());
        Assert.assertTrue(keyFile.exists());
    }
}
