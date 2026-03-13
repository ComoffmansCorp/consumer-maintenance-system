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
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class PhotoService {

    private final PhotoRepository photoRepository;
    private final InspectionActRepository inspectionActRepository;
    private final ReplacementActRepository replacementActRepository;
    private final FileStorageService fileStorageService;
    private final PhotoMapper photoMapper;
    private final CurrentTenantService currentTenantService;

    public PhotoResponse upload(PhotoUploadRequest request, MultipartFile file, Authentication authentication) {
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Uploading photo tenant={} inspectionActId={} replacementActId={} fileName={}",
                tenant.getCode(), request.getInspectionActId(), request.getReplacementActId(), file.getOriginalFilename());
        if (request.getInspectionActId() == null && request.getReplacementActId() == null) {
            throw new BadRequestException("Act is required");
        }
        var filename = fileStorageService.store(file);
        Photo photo = Photo.builder()
                .filename(filename)
                .note(request.getNote())
                .tenant(tenant)
                .build();
        if (request.getInspectionActId() != null) {
            var act = inspectionActRepository.findByIdAndTenant_Id(request.getInspectionActId(), tenant.getId())
                    .orElseThrow(() -> new ResourceNotFoundException("Inspection act not found"));
            photo.setInspectionAct(act);
        }
        if (request.getReplacementActId() != null) {
            var act = replacementActRepository.findByIdAndTenant_Id(request.getReplacementActId(), tenant.getId())
                    .orElseThrow(() -> new ResourceNotFoundException("Replacement act not found"));
            photo.setReplacementAct(act);
        }
        photoRepository.save(photo);
        log.info("Photo saved id={} tenant={}", photo.getId(), tenant.getCode());
        return photoMapper.toResponse(photo);
    }

    public List<PhotoResponse> listForInspection(Long actId, Authentication authentication) {
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Loading inspection photos actId={} tenant={}", actId, tenant.getCode());
        return photoMapper.toResponses(photoRepository.findAllByInspectionAct_IdAndTenant_Id(actId, tenant.getId()));
    }

    public List<PhotoResponse> listForReplacement(Long actId, Authentication authentication) {
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Loading replacement photos actId={} tenant={}", actId, tenant.getCode());
        return photoMapper.toResponses(photoRepository.findAllByReplacementAct_IdAndTenant_Id(actId, tenant.getId()));
    }

    public byte[] load(Long photoId, Authentication authentication) {
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Loading photo bytes photoId={} tenant={}", photoId, tenant.getCode());
        var photo = photoRepository.findByIdAndTenant_Id(photoId, tenant.getId())
                .orElseThrow(() -> new ResourceNotFoundException("Photo not found"));
        return fileStorageService.load(photo.getFilename());
    }

    public void delete(Long photoId, Authentication authentication) {
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Deleting photo photoId={} tenant={}", photoId, tenant.getCode());
        var photo = photoRepository.findByIdAndTenant_Id(photoId, tenant.getId())
                .orElseThrow(() -> new ResourceNotFoundException("Photo not found"));
        photoRepository.delete(photo);
    }
}
