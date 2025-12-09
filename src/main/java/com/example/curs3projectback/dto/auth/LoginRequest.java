package com.example.curs3projectback.dto.auth;

import lombok.Data;

@Data
public class LoginRequest {
    private String username;
    private String password;
}

