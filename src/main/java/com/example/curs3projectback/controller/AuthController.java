package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.auth.AuthResponse;
import com.example.curs3projectback.dto.auth.BootstrapSuperAdminRequest;
import com.example.curs3projectback.dto.auth.CompanyRegistrationRequest;
import com.example.curs3projectback.dto.auth.LoginRequest;
import com.example.curs3projectback.service.AuthService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@Slf4j
@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;

    @PostMapping("/register-company")
    public ResponseEntity<AuthResponse> registerCompany(@Valid @RequestBody CompanyRegistrationRequest request) {
        log.info("HTTP POST /api/auth/register-company tenantCode={} adminUsername={}", request.getTenantCode(), request.getUsername());
        return ResponseEntity.ok(authService.registerCompany(request));
    }

    @PostMapping("/bootstrap-super-admin")
    public ResponseEntity<AuthResponse> bootstrapSuperAdmin(@Valid @RequestBody BootstrapSuperAdminRequest request) {
        log.info("HTTP POST /api/auth/bootstrap-super-admin username={}", request.getUsername());
        return ResponseEntity.ok(authService.bootstrapSuperAdmin(request));
    }

    @PostMapping("/login")
    public ResponseEntity<AuthResponse> login(@Valid @RequestBody LoginRequest request) {
        log.info("HTTP POST /api/auth/login tenantCode={} username={}", request.getTenantCode(), request.getUsername());
        return ResponseEntity.ok(authService.login(request));
    }
}
