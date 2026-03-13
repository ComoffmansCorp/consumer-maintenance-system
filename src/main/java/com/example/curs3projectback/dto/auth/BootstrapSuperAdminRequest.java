package com.example.curs3projectback.dto.auth;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class BootstrapSuperAdminRequest {
    @NotBlank(message = "username is required")
    @Size(max = 100, message = "username is too long")
    private String username;

    @NotBlank(message = "password is required")
    @Size(min = 8, max = 128, message = "password length must be between 8 and 128")
    private String password;

    @NotBlank(message = "fullName is required")
    @Size(max = 255, message = "fullName is too long")
    private String fullName;
}
