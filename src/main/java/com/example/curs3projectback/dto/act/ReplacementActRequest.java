package com.example.curs3projectback.dto.act;

import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

import java.time.LocalDate;

@Data
public class ReplacementActRequest {
    @NotNull(message = "taskId is required")
    private Long taskId;
    @NotNull(message = "addressId is required")
    private Long addressId;
    @NotBlank(message = "accountNumber is required")
    @Size(max = 100, message = "accountNumber is too long")
    private String accountNumber;
    private LocalDate installationDate;
    @Valid
    @NotNull(message = "oldMeter is required")
    private MeterInfoRequest oldMeter;
    @Valid
    @NotNull(message = "newMeter is required")
    private MeterInfoRequest newMeter;
}

