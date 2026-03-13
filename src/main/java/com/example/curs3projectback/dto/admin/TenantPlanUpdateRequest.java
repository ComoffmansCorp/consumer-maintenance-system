package com.example.curs3projectback.dto.admin;

import com.example.curs3projectback.model.enums.TenantPlan;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

@Data
public class TenantPlanUpdateRequest {
    @NotNull(message = "plan is required")
    private TenantPlan plan;
}
