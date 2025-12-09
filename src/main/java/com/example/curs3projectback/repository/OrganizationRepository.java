package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.Organization;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface OrganizationRepository extends JpaRepository<Organization, Long> {
    List<Organization> findAllByNameContainingIgnoreCaseOrderByName(String query);
}

