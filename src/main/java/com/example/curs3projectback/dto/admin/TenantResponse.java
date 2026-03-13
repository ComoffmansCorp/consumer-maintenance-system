package com.example.curs3projectback.dto.admin;

import com.example.curs3projectback.model.enums.TenantPlan;
import lombok.Builder;
import lombok.Value;

import java.time.Instant;

@Value
@Builder
public class TenantResponse {
    Long id;
    String name;
    String code;
    TenantPlan plan;
    boolean active;
    Instant createdAt;
    int userLimit;
}
