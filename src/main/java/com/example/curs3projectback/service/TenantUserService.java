package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.admin.TenantUserCreateRequest;
import com.example.curs3projectback.dto.admin.TenantUserResponse;
import com.example.curs3projectback.exception.BadRequestException;
import com.example.curs3projectback.exception.ForbiddenException;
import com.example.curs3projectback.model.Tenant;
import com.example.curs3projectback.model.User;
import com.example.curs3projectback.model.enums.UserRole;
import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class TenantUserService {

    private final CurrentTenantService currentTenantService;
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    @Transactional(readOnly = true)
    public List<TenantUserResponse> findAll(Authentication authentication) {
        Tenant tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Loading tenant users tenant={}", tenant.getCode());
        return userRepository.findAllByTenant_IdOrderByFullNameAsc(tenant.getId())
                .stream()
                .map(this::toResponse)
                .toList();
    }

    @Transactional
    public TenantUserResponse create(TenantUserCreateRequest request, Authentication authentication) {
        Tenant tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Creating tenant user username={} role={} tenant={}", request.getUsername(), request.getRole(), tenant.getCode());
        if (request.getRole() == UserRole.SUPER_ADMIN) {
            throw new ForbiddenException("SUPER_ADMIN cannot be created inside tenant");
        }
        String username = request.getUsername().trim();
        if (userRepository.existsByTenant_IdAndUsernameIgnoreCase(tenant.getId(), username)) {
            throw new BadRequestException("Username is already in use in this tenant");
        }
        long currentUsers = userRepository.countByTenant_Id(tenant.getId());
        if (!tenant.getPlan().isUnlimited() && currentUsers >= tenant.getPlan().getUserLimit()) {
            throw new BadRequestException("Tenant user limit exceeded for current plan");
        }

        User user = userRepository.save(User.builder()
                .username(username)
                .password(passwordEncoder.encode(request.getPassword()))
                .fullName(request.getFullName().trim())
                .role(request.getRole())
                .tenant(tenant)
                .build());
        log.info("Created tenant user id={} username={} tenant={}", user.getId(), user.getUsername(), tenant.getCode());
        return toResponse(user);
    }

    private TenantUserResponse toResponse(User user) {
        return TenantUserResponse.builder()
                .id(user.getId())
                .username(user.getUsername())
                .fullName(user.getFullName())
                .role(user.getRole())
                .tenantCode(user.getTenant() != null ? user.getTenant().getCode() : null)
                .build();
    }
}
