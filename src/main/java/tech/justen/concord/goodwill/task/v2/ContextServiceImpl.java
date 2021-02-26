// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v2;

import com.walmartlabs.concord.runtime.v2.sdk.Context;
import com.walmartlabs.concord.runtime.v2.sdk.Variables;
import tech.justen.concord.goodwill.ContextService;

import java.util.HashSet;
import java.util.Set;

class ContextServiceImpl implements ContextService {

    private final Context ctx;

    private final Variables in;

    public ContextServiceImpl(Context ctx, Variables in) {
        this.ctx = ctx;
        this.in = in;
    }

    @Override
    public Set<String> getVariableNames() {
        Set<String> names = new HashSet<>();
        names.addAll(ctx.variables().toMap().keySet());
        names.addAll(in.toMap().keySet());
        return names;
    }

    @Override
    public Object getVariable(String name) {
        Object v = in.get(name);
        if (v != null) {
            return v;
        }
        return ctx.variables().get(name);
    }

    @Override
    public void setVariable(String name, Object value) {
        ctx.variables().set(name, value);
    }

    @Override
    public <T> T evaluate(String expr, Class<T> type) {
        return ctx.eval(expr, type);
    }

    @Override
    public <T> T getSetting(String name, T defaultValue, Class<T> clazz) {
        T value = in.get(name, null, clazz);
        if (value != null) {
            return value;
        }
        return ctx.defaultVariables().get(name, defaultValue, clazz);
    }
}
