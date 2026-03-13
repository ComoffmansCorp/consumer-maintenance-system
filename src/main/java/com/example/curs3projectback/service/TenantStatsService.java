package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.stats.TenantStatsResponse;
import com.example.curs3projectback.model.Tenant;
import com.example.curs3projectback.model.enums.TaskStatus;
import com.example.curs3projectback.repository.InspectionActRepository;
import com.example.curs3projectback.repository.ReplacementActRepository;
import com.example.curs3projectback.repository.TaskRepository;
import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Slf4j
@Service
@RequiredArgsConstructor
public class TenantStatsService {

    private static final double INSPECTION_REVENUE = 500.0;
    private static final double REPLACEMENT_REVENUE = 1500.0;

    private final CurrentTenantService currentTenantService;
    private final UserRepository userRepository;
    private final TaskRepository taskRepository;
    private final InspectionActRepository inspectionActRepository;
    private final ReplacementActRepository replacementActRepository;

    @Transactional(readOnly = true)
    public TenantStatsResponse getOverview(Authentication authentication) {
        Tenant tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Building tenant stats tenant={}", tenant.getCode());
        long activeUsers = userRepository.countByTenant_Id(tenant.getId());
        long totalTasks = taskRepository.countByTenant_Id(tenant.getId());
        long completedTasks = taskRepository.countByTenant_IdAndStatus(tenant.getId(), TaskStatus.COMPLETED);
        long inProgressTasks = taskRepository.countByTenant_IdAndStatus(tenant.getId(), TaskStatus.IN_PROGRESS);
        long inspectionActs = inspectionActRepository.countByTenant_Id(tenant.getId());
        long replacementActs = replacementActRepository.countByTenant_Id(tenant.getId());
        double completionRate = totalTasks == 0 ? 0.0 : completedTasks * 100.0 / totalTasks;
        double estimatedRevenue = inspectionActs * INSPECTION_REVENUE + replacementActs * REPLACEMENT_REVENUE;

        return TenantStatsResponse.builder()
                .tenantName(tenant.getName())
                .tenantCode(tenant.getCode())
                .plan(tenant.getPlan())
                .activeUsers(activeUsers)
                .userLimit(tenant.getPlan().getUserLimit())
                .totalTasks(totalTasks)
                .completedTasks(completedTasks)
                .inProgressTasks(inProgressTasks)
                .inspectionActs(inspectionActs)
                .replacementActs(replacementActs)
                .completionRate(completionRate)
                .estimatedRevenue(estimatedRevenue)
                .build();
    }
}
