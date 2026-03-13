package com.example.curs3projectback.dto.auth;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class LoginRequest {
    @Size(max = 100, message = "tenantCode is too long")
    private String tenantCode;

    @NotBlank(message = "username is required")
    @Size(max = 100, message = "username is too long")
    private String username;

    @NotBlank(message = "password is required")
    @Size(min = 6, max = 128, message = "password length must be between 6 and 128")
    private String password;
}

