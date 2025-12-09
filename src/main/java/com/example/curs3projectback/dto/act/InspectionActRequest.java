package com.example.curs3projectback.dto.act;

import com.example.curs3projectback.model.enums.InspectionType;
import lombok.Data;

import java.time.LocalDate;
import java.util.List;

@Data
public class InspectionActRequest {
    private Long taskId;
    private Long addressId;
    private Long consumerId;
    private LocalDate inspectionDate;
    private InspectionType inspectionType;
    private String notes;
    private List<MeterRequest> meters;
}

