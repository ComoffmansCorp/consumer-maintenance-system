package com.example.curs3projectback.model.enums;

public enum TenantPlan {
    FREE(5),
    BUSINESS(50),
    ENTERPRISE(Integer.MAX_VALUE);

    private final int userLimit;

    TenantPlan(int userLimit) {
        this.userLimit = userLimit;
    }

    public int getUserLimit() {
        return userLimit;
    }

    public boolean isUnlimited() {
        return userLimit == Integer.MAX_VALUE;
    }
}
