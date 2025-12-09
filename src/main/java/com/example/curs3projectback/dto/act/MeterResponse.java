package com.example.curs3projectback.dto.act;

import com.example.curs3projectback.model.enums.MeterType;
import lombok.Builder;
import lombok.Value;

import java.time.LocalDate;

@Value
@Builder
public class MeterResponse {
    Long id;
    MeterType type;
    String serialNumber;
    Integer manufactureYear;
    LocalDate verificationDate;
    String sealState;
    Integer transformationRatio;
}

