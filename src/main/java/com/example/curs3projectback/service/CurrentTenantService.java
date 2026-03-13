package com.example.curs3projectback.service;

import com.example.curs3projectback.exception.ForbiddenException;
import com.example.curs3projectback.model.Tenant;
import com.example.curs3projectback.model.User;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class CurrentTenantService {

    private final CurrentUserService currentUserService;

    public Tenant getRequiredTenant(Authentication authentication) {
        User user = currentUserService.getRequiredUser(authentication);
        if (user.getTenant() == null) {
            throw new ForbiddenException("Tenant context is required");
        }
        return user.getTenant();
    }
}
