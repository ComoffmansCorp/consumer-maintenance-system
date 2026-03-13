package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.admin.TenantCreateRequest;
import com.example.curs3projectback.dto.admin.TenantPlanUpdateRequest;
import com.example.curs3projectback.dto.admin.TenantResponse;
import com.example.curs3projectback.exception.BadRequestException;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.model.Tenant;
import com.example.curs3projectback.repository.TenantRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class PlatformTenantService {

    private final TenantRepository tenantRepository;

    @Transactional(readOnly = true)
    public List<TenantResponse> findAll() {
        log.info("Loading tenant list");
        return tenantRepository.findAll().stream().map(this::toResponse).toList();
    }

    @Transactional
    public TenantResponse create(TenantCreateRequest request) {
        log.info("Creating tenant code={} plan={}", request.getCode(), request.getPlan());
        validateTenantUniqueness(request.getName(), request.getCode());
        Tenant tenant = tenantRepository.save(Tenant.builder()
                .name(request.getName().trim())
                .code(normalizeTenantCode(request.getCode()))
                .plan(request.getPlan())
                .active(true)
                .build());
        return toResponse(tenant);
    }

    @Transactional
    public TenantResponse updatePlan(Long tenantId, TenantPlanUpdateRequest request) {
        log.info("Updating tenant id={} plan={}", tenantId, request.getPlan());
        Tenant tenant = getRequiredTenant(tenantId);
        tenant.setPlan(request.getPlan());
        return toResponse(tenantRepository.save(tenant));
    }

    @Transactional
    public TenantResponse deactivate(Long tenantId) {
        log.info("Deactivating tenant id={}", tenantId);
        Tenant tenant = getRequiredTenant(tenantId);
        tenant.setActive(false);
        return toResponse(tenantRepository.save(tenant));
    }

    public String normalizeTenantCode(String code) {
        return code.trim().toLowerCase();
    }

    private Tenant getRequiredTenant(Long tenantId) {
        return tenantRepository.findById(tenantId)
                .orElseThrow(() -> new ResourceNotFoundException("Tenant not found"));
    }

    private void validateTenantUniqueness(String name, String code) {
        if (tenantRepository.existsByNameIgnoreCase(name.trim())) {
            throw new BadRequestException("Tenant name is already in use");
        }
        if (tenantRepository.existsByCodeIgnoreCase(normalizeTenantCode(code))) {
            throw new BadRequestException("Tenant code is already in use");
        }
    }

    private TenantResponse toResponse(Tenant tenant) {
        return TenantResponse.builder()
                .id(tenant.getId())
                .name(tenant.getName())
                .code(tenant.getCode())
                .plan(tenant.getPlan())
                .active(tenant.isActive())
                .createdAt(tenant.getCreatedAt())
                .userLimit(tenant.getPlan().getUserLimit())
                .build();
    }
}
