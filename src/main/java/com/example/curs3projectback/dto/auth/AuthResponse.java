package com.example.curs3projectback.dto.auth;

import com.example.curs3projectback.model.enums.UserRole;
import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class AuthResponse {
    private String token;
    private String fullName;
    private UserRole role;
}

