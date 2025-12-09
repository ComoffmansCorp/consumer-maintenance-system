package com.example.curs3projectback.mapper;

import com.example.curs3projectback.dto.photo.PhotoResponse;
import com.example.curs3projectback.model.Photo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring")
public interface PhotoMapper {

    @Mapping(target = "url", expression = "java(\"/api/photos/\" + photo.getId())")
    PhotoResponse toResponse(Photo photo);

    List<PhotoResponse> toResponses(List<Photo> photos);
}

