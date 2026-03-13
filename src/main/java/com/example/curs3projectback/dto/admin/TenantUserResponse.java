package com.example.curs3projectback.dto.admin;

import com.example.curs3projectback.model.enums.UserRole;
import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class TenantUserResponse {
    Long id;
    String username;
    String fullName;
    UserRole role;
    String tenantCode;
}
