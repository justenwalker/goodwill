// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import java.util.Set;

public interface ContextService {

  Set<String> getVariableNames();

  Object getVariable(String name);

  void setVariable(String name, Object value);

  <T> T getSetting(String name, T defaultValue, Class<T> type);

  <T> T evaluate(String expr, Class<T> type);
}
