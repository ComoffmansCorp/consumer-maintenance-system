package com.example.curs3projectback.dto.stats;

import com.example.curs3projectback.model.enums.TenantPlan;
import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class TenantStatsResponse {
    String tenantName;
    String tenantCode;
    TenantPlan plan;
    long activeUsers;
    int userLimit;
    long totalTasks;
    long completedTasks;
    long inProgressTasks;
    long inspectionActs;
    long replacementActs;
    double completionRate;
    double estimatedRevenue;
}
