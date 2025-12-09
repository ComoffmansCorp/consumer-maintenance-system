package com.example.curs3projectback.dto.auth;

import com.example.curs3projectback.model.enums.UserRole;
import lombok.Data;

@Data
public class RegisterRequest {
    private String username;
    private String password;
    private String fullName;
    private UserRole role;
}

