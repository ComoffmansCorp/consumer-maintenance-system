package com.example.curs3projectback.dto.admin;

import com.example.curs3projectback.model.enums.UserRole;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class TenantUserCreateRequest {
    @NotBlank(message = "username is required")
    @Size(max = 100, message = "username is too long")
    private String username;

    @NotBlank(message = "password is required")
    @Size(min = 6, max = 128, message = "password length must be between 6 and 128")
    private String password;

    @NotBlank(message = "fullName is required")
    @Size(max = 255, message = "fullName is too long")
    private String fullName;

    @NotNull(message = "role is required")
    private UserRole role;
}
