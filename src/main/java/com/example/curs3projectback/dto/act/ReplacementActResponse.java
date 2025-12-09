package com.example.curs3projectback.dto.act;

import lombok.Builder;
import lombok.Value;

import java.time.LocalDate;

@Value
@Builder
public class ReplacementActResponse {
    Long id;
    Long taskId;
    String address;
    String accountNumber;
    LocalDate installationDate;
    MeterInfoResponse oldMeter;
    MeterInfoResponse newMeter;
    int photosCount;
}

