package com.example.curs3projectback.dto.admin;

import com.example.curs3projectback.model.enums.TenantPlan;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class TenantCreateRequest {
    @NotBlank(message = "name is required")
    @Size(max = 150, message = "name is too long")
    private String name;

    @NotBlank(message = "code is required")
    @Size(max = 100, message = "code is too long")
    private String code;

    @NotNull(message = "plan is required")
    private TenantPlan plan;
}
