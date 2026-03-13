package com.example.curs3projectback.dto.auth;

import com.example.curs3projectback.model.enums.TenantPlan;
import com.example.curs3projectback.model.enums.UserRole;
import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class AuthResponse {
    String token;
    String fullName;
    UserRole role;
    Long tenantId;
    String tenantCode;
    String tenantName;
    TenantPlan tenantPlan;
}

