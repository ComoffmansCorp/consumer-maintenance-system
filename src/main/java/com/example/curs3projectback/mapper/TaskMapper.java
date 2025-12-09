package com.example.curs3projectback.mapper;

import com.example.curs3projectback.dto.task.TaskResponse;
import com.example.curs3projectback.model.Task;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring")
public interface TaskMapper {

    @Mapping(target = "address", expression = "java(formatAddress(task))")
    @Mapping(target = "consumerName", expression = "java(task.getAddress() != null && task.getAddress().getConsumer() != null ? task.getAddress().getConsumer().getName() : null)")
    TaskResponse toResponse(Task task);

    List<TaskResponse> toResponses(List<Task> tasks);

    default String formatAddress(Task task) {
        if (task.getAddress() == null) {
            return null;
        }
        var address = task.getAddress();
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

