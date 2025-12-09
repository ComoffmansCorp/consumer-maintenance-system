package com.example.curs3projectback.dto.act;

import com.example.curs3projectback.model.enums.MeterType;
import lombok.Data;

import java.time.LocalDate;

@Data
public class MeterRequest {
    private MeterType type;
    private String serialNumber;
    private Integer manufactureYear;
    private LocalDate verificationDate;
    private String sealState;
    private Integer transformationRatio;
}

