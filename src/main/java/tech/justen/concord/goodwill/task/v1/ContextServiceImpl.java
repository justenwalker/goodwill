// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v1;

import com.walmartlabs.concord.sdk.Context;
import tech.justen.concord.goodwill.ContextService;

import java.util.Map;
import java.util.Set;

class ContextServiceImpl implements ContextService {

    private final Context ctx;

    public ContextServiceImpl(Context ctx) {
        this.ctx = ctx;
    }

    @Override
    public Set<String> getVariableNames() {
        return ctx.toMap().keySet();
    }

    @Override
    public Object getVariable(String name) {
        return ctx.getVariable(name);
    }

    @Override
    public void setVariable(String name, Object value) {
        ctx.setVariable(name, value);
    }

    @Override
    public <T> T evaluate(String expr, Class<T> type) {
        return ctx.eval(expr, type);
    }

    @Override
    @SuppressWarnings("unchecked")
    public <T> T getSetting(String name, T defaultValue, Class<T> clazz) {
        Object value = ctx.getVariable(name);
        if (value != null) {
            return clazz.cast(value);
        }
        Map<String, Object> defaults = (Map<String, Object>) ctx.getVariable("goodwillDefaults");
        if (defaults == null) {
            return defaultValue;
        }
        value = defaults.getOrDefault(name, defaultValue);
        return clazz.cast(value);
    }
}
