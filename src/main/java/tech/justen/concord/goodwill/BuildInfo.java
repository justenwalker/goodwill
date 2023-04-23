/*
 * Copyright 2021, Justen Walker and the goodwill contributors
 * SPDX-License-Identifier: Apache-2.0
 */

package tech.justen.concord.goodwill;

import java.io.IOException;
import java.util.Properties;

public enum BuildInfo {
    INSTANCE;

    private Properties properties;

    BuildInfo() {
        properties = new Properties();
        try {
            properties.load(getClass().getResourceAsStream("/version.properties"));
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static String getProperty(String key) {
        return INSTANCE.properties.getProperty(key);
    }

    public static String getVersion() {
        return getProperty("build.version");
    }
    public static String getGitCommit() {
        return getProperty("build.gitCommit");
    }
    public static String getBuildTimestamp() {
        return getProperty("build.timestamp");
    }
}
