// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import java.util.Map;

public interface SecretService {

    String PUBLIC_KEY = "public";

    String PRIVATE_KEY = "private";

    String USERNAME = "username";

    String PASSWORD = "password";

    String exportAsString(String orgName,
                          String name,
                          String password) throws Exception;


    Map<String, String> exportKeyAsFile(String orgName,
                                        String name,
                                        String password) throws Exception;

    Map<String, String> exportCredentials(String orgName,
                                          String name,
                                          String password) throws Exception;

    String exportAsFile(String orgName,
                        String name,
                        String password) throws Exception;

    String decryptString(String s) throws Exception;

    String encryptString(String orgName,
                         String projectName,
                         String value) throws Exception;
}
