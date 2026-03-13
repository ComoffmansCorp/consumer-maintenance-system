package com.example.curs3projectback.dto.act;

import com.example.curs3projectback.model.enums.InspectionType;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

import java.time.LocalDate;
import java.util.List;

@Data
public class InspectionActRequest {
    @NotNull(message = "taskId is required")
    private Long taskId;
    @NotNull(message = "addressId is required")
    private Long addressId;
    private Long consumerId;
    private LocalDate inspectionDate;
    @NotNull(message = "inspectionType is required")
    private InspectionType inspectionType;
    @Size(max = 2000, message = "notes is too long")
    private String notes;
    @Valid
    private List<MeterRequest> meters;
}

