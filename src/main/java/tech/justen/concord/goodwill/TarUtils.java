// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream;
import org.apache.commons.compress.utils.IOUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.PosixFilePermission;
import java.nio.file.attribute.PosixFilePermissions;
import java.util.Set;
import java.util.zip.GZIPInputStream;

public class TarUtils {

    private static final Logger log = LoggerFactory.getLogger(TarUtils.class);

    public static void extractTarball(Path tarBall, Path out) throws IOException {
        TarArchiveInputStream is = new TarArchiveInputStream(new GZIPInputStream(new FileInputStream(tarBall.toFile())));
        try {
            TarArchiveEntry entry;
            while ((entry = (TarArchiveEntry) is.getNextEntry()) != null) {
                final File outputFile = new File(out.toFile(), entry.getName());
                if (entry.isDirectory()) {
                    if (!outputFile.exists()) {
                        outputFile.mkdirs();
                    }
                } else {
                    outputFile.getParentFile().mkdirs();
                    log.debug("tar.gz extract {} => {}", entry.getName(), outputFile);
                    try (OutputStream outputFileStream = new FileOutputStream(outputFile)) {
                        IOUtils.copy(is, outputFileStream);
                    }
                    setPosixFilePermissions(outputFile.toPath(), entry.getMode());
                }
            }
        } finally {
            is.close();
        }
    }


    public static void setPosixFilePermissions(Path path, int mode) throws IOException {
        final char[] ss = {'r', 'w', 'x', 'r', 'w', 'x', 'r', 'w', 'x'};
        int i = ss.length - 1;
        for (int b = 1; b < 512; b <<= 1) {
            if ((b & mode) == 0) {
                ss[i] = '-';
            }
            i--;
        }
        String sperms = new String(ss);
        //System.out.printf("0%o -> %s %s\n", mode, sperms, path);
        Files.setPosixFilePermissions(path, PosixFilePermissions.fromString(sperms));
    }
}