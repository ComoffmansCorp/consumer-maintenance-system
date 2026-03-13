package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.admin.TenantUserCreateRequest;
import com.example.curs3projectback.dto.admin.TenantUserResponse;
import com.example.curs3projectback.service.TenantUserService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/api/tenant/users")
@RequiredArgsConstructor
@PreAuthorize("hasRole('TENANT_ADMIN')")
public class TenantUserController {

    private final TenantUserService tenantUserService;

    @GetMapping
    public ResponseEntity<List<TenantUserResponse>> findAll(Authentication authentication) {
        log.info("HTTP GET /api/tenant/users userId={}", authentication.getName());
        return ResponseEntity.ok(tenantUserService.findAll(authentication));
    }

    @PostMapping
    public ResponseEntity<TenantUserResponse> create(@Valid @RequestBody TenantUserCreateRequest request,
                                                     Authentication authentication) {
        log.info("HTTP POST /api/tenant/users userId={} username={} role={}",
                authentication.getName(), request.getUsername(), request.getRole());
        return ResponseEntity.ok(tenantUserService.create(request, authentication));
    }
}
