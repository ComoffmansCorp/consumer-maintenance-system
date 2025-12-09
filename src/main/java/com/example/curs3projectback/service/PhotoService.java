package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.photo.PhotoResponse;
import com.example.curs3projectback.dto.photo.PhotoUploadRequest;
import com.example.curs3projectback.exception.BadRequestException;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.mapper.PhotoMapper;
import com.example.curs3projectback.model.Photo;
import com.example.curs3projectback.repository.InspectionActRepository;
import com.example.curs3projectback.repository.PhotoRepository;
import com.example.curs3projectback.repository.ReplacementActRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

@Service
@RequiredArgsConstructor
public class PhotoService {

    private final PhotoRepository photoRepository;
    private final InspectionActRepository inspectionActRepository;
    private final ReplacementActRepository replacementActRepository;
    private final FileStorageService fileStorageService;
    private final PhotoMapper photoMapper;

    public PhotoResponse upload(PhotoUploadRequest request, MultipartFile file) {
        if (request.getInspectionActId() == null && request.getReplacementActId() == null) {
            throw new BadRequestException("Не указан акт");
        }
        var filename = fileStorageService.store(file);
        Photo photo = Photo.builder()
                .filename(filename)
                .note(request.getNote())
                .build();
        if (request.getInspectionActId() != null) {
            var act = inspectionActRepository.findById(request.getInspectionActId())
                    .orElseThrow(() -> new ResourceNotFoundException("Акт осмотра не найден"));
            photo.setInspectionAct(act);
        }
        if (request.getReplacementActId() != null) {
            var act = replacementActRepository.findById(request.getReplacementActId())
                    .orElseThrow(() -> new ResourceNotFoundException("Акт замены не найден"));
            photo.setReplacementAct(act);
        }
        photoRepository.save(photo);
        return photoMapper.toResponse(photo);
    }

    public List<PhotoResponse> listForInspection(Long actId) {
        return photoMapper.toResponses(photoRepository.findAllByInspectionAct_Id(actId));
    }

    public List<PhotoResponse> listForReplacement(Long actId) {
        return photoMapper.toResponses(photoRepository.findAllByReplacementAct_Id(actId));
    }

    public byte[] load(Long photoId) {
        var photo = photoRepository.findById(photoId)
                .orElseThrow(() -> new ResourceNotFoundException("Фото не найдено"));
        return fileStorageService.load(photo.getFilename());
    }

    public void delete(Long photoId) {
        var photo = photoRepository.findById(photoId)
                .orElseThrow(() -> new ResourceNotFoundException("Фото не найдено"));
        photoRepository.delete(photo);
    }
}

