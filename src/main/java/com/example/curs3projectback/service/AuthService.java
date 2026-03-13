package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.auth.AuthResponse;
import com.example.curs3projectback.dto.auth.BootstrapSuperAdminRequest;
import com.example.curs3projectback.dto.auth.CompanyRegistrationRequest;
import com.example.curs3projectback.dto.auth.LoginRequest;
import com.example.curs3projectback.exception.BadRequestException;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.model.Tenant;
import com.example.curs3projectback.model.User;
import com.example.curs3projectback.model.enums.TenantPlan;
import com.example.curs3projectback.model.enums.UserRole;
import com.example.curs3projectback.repository.TenantRepository;
import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Slf4j
@Service
@RequiredArgsConstructor
public class AuthService {

    private static final String LOGIN_SEPARATOR = "::";

    private final UserRepository userRepository;
    private final TenantRepository tenantRepository;
    private final PasswordEncoder passwordEncoder;
    private final AuthenticationManager authenticationManager;
    private final JwtService jwtService;
    private final PlatformTenantService platformTenantService;

    @Transactional
    public AuthResponse registerCompany(CompanyRegistrationRequest request) {
        String tenantCode = platformTenantService.normalizeTenantCode(request.getTenantCode());
        log.info("Registering tenant tenantCode={} tenantAdmin={}", tenantCode, request.getUsername());
        if (tenantRepository.existsByNameIgnoreCase(request.getTenantName().trim())) {
            throw new BadRequestException("Tenant name is already in use");
        }
        if (tenantRepository.existsByCodeIgnoreCase(tenantCode)) {
            throw new BadRequestException("Tenant code is already in use");
        }

        Tenant tenant = tenantRepository.save(Tenant.builder()
                .name(request.getTenantName().trim())
                .code(tenantCode)
                .plan(request.getPlan() == null ? TenantPlan.FREE : request.getPlan())
                .active(true)
                .build());

        User tenantAdmin = createUser(
                tenant,
                request.getUsername(),
                request.getPassword(),
                request.getFullName(),
                UserRole.TENANT_ADMIN
        );
        return buildAuthResponse(tenantAdmin);
    }

    @Transactional
    public AuthResponse bootstrapSuperAdmin(BootstrapSuperAdminRequest request) {
        log.info("Bootstrapping SUPER_ADMIN username={}", request.getUsername());
        if (userRepository.countByRole(UserRole.SUPER_ADMIN) > 0) {
            throw new BadRequestException("SUPER_ADMIN already exists");
        }
        User superAdmin = createUser(null, request.getUsername(), request.getPassword(), request.getFullName(), UserRole.SUPER_ADMIN);
        return buildAuthResponse(superAdmin);
    }

    public AuthResponse login(LoginRequest request) {
        String loginKey = buildLoginKey(request);
        log.info("Authenticating username={} tenantCode={}", request.getUsername(), request.getTenantCode());
        authenticationManager.authenticate(new UsernamePasswordAuthenticationToken(loginKey, request.getPassword()));
        User user = resolveUserForLogin(request);
        return buildAuthResponse(user);
    }

    private User resolveUserForLogin(LoginRequest request) {
        String username = request.getUsername().trim();
        String tenantCode = normalizeOptionalTenantCode(request.getTenantCode());
        if (tenantCode == null) {
            return userRepository.findByTenantIsNullAndUsernameIgnoreCase(username)
                    .orElseThrow(() -> new BadRequestException("Invalid credentials"));
        }
        return userRepository.findByTenant_CodeIgnoreCaseAndUsernameIgnoreCase(tenantCode, username)
                .orElseThrow(() -> new BadRequestException("Invalid credentials"));
    }

    private User createUser(Tenant tenant, String username, String password, String fullName, UserRole role) {
        String normalizedUsername = username.trim();
        if (tenant != null) {
            if (userRepository.existsByTenant_IdAndUsernameIgnoreCase(tenant.getId(), normalizedUsername)) {
                throw new BadRequestException("Username is already in use in this tenant");
            }
            long currentUsers = userRepository.countByTenant_Id(tenant.getId());
            if (!tenant.getPlan().isUnlimited() && currentUsers >= tenant.getPlan().getUserLimit()) {
                throw new BadRequestException("Tenant user limit exceeded for current plan");
            }
        } else if (userRepository.findByTenantIsNullAndUsernameIgnoreCase(normalizedUsername).isPresent()) {
            throw new BadRequestException("Username is already in use");
        }

        User savedUser = userRepository.save(User.builder()
                .username(normalizedUsername)
                .password(passwordEncoder.encode(password))
                .fullName(fullName.trim())
                .role(role)
                .tenant(tenant)
                .build());
        log.info("Created user id={} username={} role={} tenant={}",
                savedUser.getId(),
                savedUser.getUsername(),
                savedUser.getRole(),
                savedUser.getTenant() != null ? savedUser.getTenant().getCode() : "platform");
        return savedUser;
    }

    private AuthResponse buildAuthResponse(User user) {
        String token = jwtService.generateToken(jwtClaims(user), toSecurityUser(user));
        return AuthResponse.builder()
                .token(token)
                .fullName(user.getFullName())
                .role(user.getRole())
                .tenantId(user.getTenant() != null ? user.getTenant().getId() : null)
                .tenantCode(user.getTenant() != null ? user.getTenant().getCode() : null)
                .tenantName(user.getTenant() != null ? user.getTenant().getName() : null)
                .tenantPlan(user.getTenant() != null ? user.getTenant().getPlan() : null)
                .build();
    }

    private Map<String, Object> jwtClaims(User user) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", user.getId());
        claims.put("role", user.getRole().name());
        if (user.getTenant() != null) {
            claims.put("tenantId", user.getTenant().getId());
            claims.put("tenantCode", user.getTenant().getCode());
        }
        return claims;
    }

    private org.springframework.security.core.userdetails.User toSecurityUser(User user) {
        return new org.springframework.security.core.userdetails.User(
                String.valueOf(user.getId()),
                user.getPassword(),
                List.of(new SimpleGrantedAuthority("ROLE_" + user.getRole().name()))
        );
    }

    private String buildLoginKey(LoginRequest request) {
        String tenantCode = normalizeOptionalTenantCode(request.getTenantCode());
        if (tenantCode == null) {
            return request.getUsername().trim();
        }
        return tenantCode + LOGIN_SEPARATOR + request.getUsername().trim();
    }

    private String normalizeOptionalTenantCode(String tenantCode) {
        if (tenantCode == null || tenantCode.isBlank()) {
            return null;
        }
        return platformTenantService.normalizeTenantCode(tenantCode);
    }
}
