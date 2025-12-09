package com.example.curs3projectback.mapper;

import com.example.curs3projectback.dto.act.InspectionActResponse;
import com.example.curs3projectback.dto.act.MeterResponse;
import com.example.curs3projectback.model.InspectionAct;
import com.example.curs3projectback.model.Meter;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring")
public interface InspectionActMapper {

    @Mapping(target = "taskId", source = "task.id")
    @Mapping(target = "address", expression = "java(formatAddress(act))")
    @Mapping(target = "consumer", expression = "java(act.getConsumer() != null ? act.getConsumer().getName() : null)")
    @Mapping(target = "meters", source = "meters")
    @Mapping(target = "photosCount", expression = "java(act.getPhotos() != null ? act.getPhotos().size() : 0)")
    InspectionActResponse toResponse(InspectionAct act);

    List<InspectionActResponse> toResponses(List<InspectionAct> acts);

    @Mapping(target = "id", source = "id")
    MeterResponse toMeterResponse(Meter meter);

    default String formatAddress(InspectionAct act) {
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

