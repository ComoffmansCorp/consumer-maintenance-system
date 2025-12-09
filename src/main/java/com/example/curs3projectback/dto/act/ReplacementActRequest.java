package com.example.curs3projectback.dto.act;

import lombok.Data;

import java.time.LocalDate;

@Data
public class ReplacementActRequest {
    private Long taskId;
    private Long addressId;
    private String accountNumber;
    private LocalDate installationDate;
    private MeterInfoRequest oldMeter;
    private MeterInfoRequest newMeter;
}

