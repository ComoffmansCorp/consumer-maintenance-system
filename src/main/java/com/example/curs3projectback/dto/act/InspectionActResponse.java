package com.example.curs3projectback.dto.act;

import com.example.curs3projectback.model.enums.InspectionType;
import lombok.Builder;
import lombok.Value;

import java.time.LocalDate;
import java.util.List;

@Value
@Builder
public class InspectionActResponse {
    Long id;
    Long taskId;
    String address;
    String consumer;
    LocalDate inspectionDate;
    InspectionType inspectionType;
    String notes;
    List<MeterResponse> meters;
    int photosCount;
}

