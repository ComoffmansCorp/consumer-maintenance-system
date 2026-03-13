package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.Organization;
import com.example.curs3projectback.model.Tenant;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface OrganizationRepository extends JpaRepository<Organization, Long> {
    List<Organization> findAllByNameContainingIgnoreCaseOrderByName(String query);
    List<Organization> findAllByTenantAndNameContainingIgnoreCaseOrderByName(Tenant tenant, String query);
    Optional<Organization> findByIdAndTenant_Id(Long id, Long tenantId);
}

