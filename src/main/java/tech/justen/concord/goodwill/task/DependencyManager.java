// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task;

import java.io.IOException;
import java.net.URI;
import java.nio.file.Path;

public interface DependencyManager {
    Path resolve(URI uri) throws IOException;
}
