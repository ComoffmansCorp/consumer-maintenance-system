package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.admin.TenantCreateRequest;
import com.example.curs3projectback.dto.admin.TenantPlanUpdateRequest;
import com.example.curs3projectback.dto.admin.TenantResponse;
import com.example.curs3projectback.service.PlatformTenantService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/api/platform/tenants")
@RequiredArgsConstructor
@PreAuthorize("hasRole('SUPER_ADMIN')")
public class PlatformTenantController {

    private final PlatformTenantService platformTenantService;

    @GetMapping
    public ResponseEntity<List<TenantResponse>> findAll() {
        log.info("HTTP GET /api/platform/tenants");
        return ResponseEntity.ok(platformTenantService.findAll());
    }

    @PostMapping
    public ResponseEntity<TenantResponse> create(@Valid @RequestBody TenantCreateRequest request) {
        log.info("HTTP POST /api/platform/tenants code={} plan={}", request.getCode(), request.getPlan());
        return ResponseEntity.ok(platformTenantService.create(request));
    }

    @PatchMapping("/{tenantId}/plan")
    public ResponseEntity<TenantResponse> updatePlan(@PathVariable Long tenantId,
                                                     @Valid @RequestBody TenantPlanUpdateRequest request) {
        log.info("HTTP PATCH /api/platform/tenants/{}/plan newPlan={}", tenantId, request.getPlan());
        return ResponseEntity.ok(platformTenantService.updatePlan(tenantId, request));
    }

    @DeleteMapping("/{tenantId}")
    public ResponseEntity<TenantResponse> deactivate(@PathVariable Long tenantId) {
        log.info("HTTP DELETE /api/platform/tenants/{}", tenantId);
        return ResponseEntity.ok(platformTenantService.deactivate(tenantId));
    }
}
