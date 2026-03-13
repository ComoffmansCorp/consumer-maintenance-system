package com.example.curs3projectback.dto.act;

import com.example.curs3projectback.model.enums.MeterType;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

import java.time.LocalDate;

@Data
public class MeterRequest {
    @NotNull(message = "type is required")
    private MeterType type;
    @NotBlank(message = "serialNumber is required")
    @Size(max = 128, message = "serialNumber is too long")
    private String serialNumber;
    @Min(value = 1950, message = "manufactureYear is too small")
    @Max(value = 2100, message = "manufactureYear is too large")
    private Integer manufactureYear;
    private LocalDate verificationDate;
    @Size(max = 255, message = "sealState is too long")
    private String sealState;
    @Min(value = 1, message = "transformationRatio must be positive")
    private Integer transformationRatio;
}

