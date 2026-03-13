package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.Photo;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface PhotoRepository extends JpaRepository<Photo, Long> {
    List<Photo> findAllByInspectionAct_IdAndTenant_Id(Long actId, Long tenantId);
    List<Photo> findAllByReplacementAct_IdAndTenant_Id(Long actId, Long tenantId);
    Optional<Photo> findByIdAndTenant_Id(Long photoId, Long tenantId);
}

