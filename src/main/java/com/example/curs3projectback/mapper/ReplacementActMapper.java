package com.example.curs3projectback.mapper;

import com.example.curs3projectback.dto.act.MeterInfoResponse;
import com.example.curs3projectback.dto.act.ReplacementActResponse;
import com.example.curs3projectback.model.MeterInfo;
import com.example.curs3projectback.model.ReplacementAct;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring")
public interface ReplacementActMapper {

    @Mapping(target = "taskId", source = "task.id")
    @Mapping(target = "address", expression = "java(formatAddress(act))")
    @Mapping(target = "oldMeter", source = "oldMeter")
    @Mapping(target = "newMeter", source = "newMeter")
    @Mapping(target = "photosCount", expression = "java(act.getPhotos() != null ? act.getPhotos().size() : 0)")
    ReplacementActResponse toResponse(ReplacementAct act);

    List<ReplacementActResponse> toResponses(List<ReplacementAct> acts);

    MeterInfoResponse toMeterInfoResponse(MeterInfo info);

    default String formatAddress(ReplacementAct act) {
        if (act.getAddress() == null) {
            return null;
        }
        var address = act.getAddress();
        StringBuilder sb = new StringBuilder();
        sb.append(address.getStreet()).append(", д. ").append(address.getHouse());
        if (address.getBuilding() != null && !address.getBuilding().isBlank()) {
            sb.append(", корп. ").append(address.getBuilding());
        }
        if (address.getApartment() != null && !address.getApartment().isBlank()) {
            sb.append(", кв. ").append(address.getApartment());
        }
        return sb.toString();
    }
}

