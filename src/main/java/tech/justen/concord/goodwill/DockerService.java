// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import java.io.IOException;

public interface DockerService {
  int start(DockerContainer container, LogCallback outCallback, LogCallback errCallback)
      throws IOException, InterruptedException;

  interface LogCallback {
    void onLog(String line);
  }
}
