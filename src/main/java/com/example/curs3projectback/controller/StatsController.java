package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.stats.TenantStatsResponse;
import com.example.curs3projectback.service.TenantStatsService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@Slf4j
@RestController
@RequestMapping("/api/stats")
@RequiredArgsConstructor
@PreAuthorize("hasAnyRole('TENANT_ADMIN','DISPATCHER','ELECTRICIAN')")
public class StatsController {

    private final TenantStatsService tenantStatsService;

    @GetMapping("/overview")
    public ResponseEntity<TenantStatsResponse> overview(Authentication authentication) {
        log.info("HTTP GET /api/stats/overview userId={}", authentication.getName());
        return ResponseEntity.ok(tenantStatsService.getOverview(authentication));
    }
}
