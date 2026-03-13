package com.example.curs3projectback.dto.auth;

import com.example.curs3projectback.model.enums.TenantPlan;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class CompanyRegistrationRequest {
    @NotBlank(message = "tenantName is required")
    @Size(max = 150, message = "tenantName is too long")
    private String tenantName;

    @NotBlank(message = "tenantCode is required")
    @Size(max = 100, message = "tenantCode is too long")
    private String tenantCode;

    private TenantPlan plan = TenantPlan.FREE;

    @NotBlank(message = "username is required")
    @Size(max = 100, message = "username is too long")
    private String username;

    @NotBlank(message = "password is required")
    @Size(min = 6, max = 128, message = "password length must be between 6 and 128")
    private String password;

    @NotBlank(message = "fullName is required")
    @Size(max = 255, message = "fullName is too long")
    private String fullName;
}
