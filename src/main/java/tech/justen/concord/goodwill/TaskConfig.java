// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import java.nio.file.Path;

public interface TaskConfig {
  String processId();

  Path workingDirectory();

  String orgName();

  String orgId();

  String projectName();

  String projectId();

  String repoName();

  String repoId();

  String repoUrl();
}
