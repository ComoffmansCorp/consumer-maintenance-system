package com.example.curs3projectback.config;

import com.example.curs3projectback.model.*;
import com.example.curs3projectback.model.enums.*;
import com.example.curs3projectback.repository.*;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDate;
import java.util.List;

/**
 * Создаёт демо-данные при первом запуске.
 * Пароль для всех: demo123  (кроме superadmin → admin123)
 *
 * Логин для Android:
 *   Код компании: esc-ural
 *   Логин: kozlov_d  |  fedorova_a  |  volkov_r
 *   Пароль: demo123
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class DemoDataInitializer implements ApplicationRunner {

    private final TenantRepository     tenantRepo;
    private final UserRepository       userRepo;
    private final OrganizationRepository orgRepo;
    private final AddressRepository    addressRepo;
    private final TaskRepository       taskRepo;
    private final PasswordEncoder      encoder;

    @Override
    @Transactional
    public void run(ApplicationArguments args) {
        if (tenantRepo.count() > 0) {
            log.info("Demo data already exists — skipping initialization");
            return;
        }

        log.info("=== Initializing demo data ===");

        String adminPass  = encoder.encode("admin123");
        String demoPass   = encoder.encode("demo123");

        // ── Superadmin ──────────────────────────────────────────────────────────
        userRepo.save(User.builder()
                .username("superadmin")
                .password(adminPass)
                .fullName("Platform Super Admin")
                .role(UserRole.SUPER_ADMIN)
                .build());

        // ── Tenants ─────────────────────────────────────────────────────────────
        Tenant esc = save(t("ЭнергоСервис Урал", "esc-ural",     TenantPlan.ENTERPRISE));
        Tenant gaz = save(t("ГазМонтаж Плюс",    "gazmontazh",   TenantPlan.BUSINESS));
        Tenant str = save(t("СтройКонтроль",     "stroykontrol", TenantPlan.BUSINESS));
        Tenant mw  = save(t("МегаВатт Сервис",   "megawatt",     TenantPlan.ENTERPRISE));
        Tenant ng  = save(t("North Grid",        "north-grid",   TenantPlan.FREE));

        // ── Users — esc-ural ────────────────────────────────────────────────────
        userRepo.save(user("ivanov_p",   adminPass, "Иванов Павел Алексеевич",     UserRole.TENANT_ADMIN, esc));
        userRepo.save(user("smirnova_o", demoPass,  "Смирнова Оксана Дмитриевна",  UserRole.DISPATCHER,   esc));
        User kozlov   = userRepo.save(user("kozlov_d",   demoPass, "Козлов Дмитрий Игоревич",     UserRole.ELECTRICIAN, esc));
        User fedorova = userRepo.save(user("fedorova_a", demoPass, "Фёдорова Анна Сергеевна",     UserRole.ELECTRICIAN, esc));
        User volkov   = userRepo.save(user("volkov_r",   demoPass, "Волков Роман Николаевич",     UserRole.ELECTRICIAN, esc));

        // ── Users — other tenants ───────────────────────────────────────────────
        userRepo.save(user("petrov_a",   adminPass, "Петров Алексей Витальевич",  UserRole.TENANT_ADMIN, gaz));
        User novikov = userRepo.save(user("novikov_s", demoPass, "Новиков Сергей Павлович",   UserRole.ELECTRICIAN, gaz));
        userRepo.save(user("morozova_e", adminPass, "Морозова Елена Владимировна",UserRole.TENANT_ADMIN, str));
        User titov = userRepo.save(user("titov_k",   demoPass, "Титов Константин Андреевич", UserRole.ELECTRICIAN, str));
        userRepo.save(user("belov_v",    adminPass, "Белов Виктор Михайлович",    UserRole.TENANT_ADMIN, mw));
        User gorbunov = userRepo.save(user("gorbunov_n", demoPass, "Горбунов Никита Олегович",   UserRole.ELECTRICIAN, mw));
        User ershov   = userRepo.save(user("ershov_p",   demoPass, "Ершов Павел Геннадьевич",    UserRole.ELECTRICIAN, mw));
        userRepo.save(user("northadmin", adminPass, "North Grid Admin",           UserRole.TENANT_ADMIN, ng));

        // ── Organizations ───────────────────────────────────────────────────────
        Organization energo = orgRepo.save(org("ОАО ЭнергоПоставка", OrganizationType.COMMERCIAL,  esc));
        Organization school = orgRepo.save(org("Школа №14",          OrganizationType.GOVERNMENT,  esc));
        Organization tc     = orgRepo.save(org("ТЦ Мегаполис",       OrganizationType.COMMERCIAL,  esc));
        Organization hosp   = orgRepo.save(org("Больница №2",        OrganizationType.GOVERNMENT,  esc));
        Organization gazt   = orgRepo.save(org("ГазТрейд ООО",       OrganizationType.COMMERCIAL,  gaz));
        Organization mkd    = orgRepo.save(org("МКД Ленина 15",       OrganizationType.RESIDENTIAL, gaz));
        Organization zhk    = orgRepo.save(org("ЖК Панорама",         OrganizationType.RESIDENTIAL, str));
        Organization zavod  = orgRepo.save(org("Завод Прогресс",      OrganizationType.COMMERCIAL,  mw));
        Organization mega   = orgRepo.save(org("МегаЭнерго",          OrganizationType.COMMERCIAL,  mw));
        Organization nschool= orgRepo.save(org("North School",        OrganizationType.GOVERNMENT,  ng));

        // ── Addresses — esc-ural ────────────────────────────────────────────────
        Address a1 = addr("ул. Ленина",       "12", null,  "34",  esc, energo);
        Address a2 = addr("ул. Ленина",       "12", null,  "7",   esc, energo);
        Address a3 = addr("пр. Мира",         "5",  "1",   "7",   esc, energo);
        Address a4 = addr("ул. Гагарина",     "3",  null,  "15",  esc, school);
        Address a5 = addr("ул. Гагарина",     "44", null,  "24",  esc, tc);
        Address a6 = addr("ул. Советская",    "88", null,  "2",   esc, hosp);
        Address a7 = addr("ул. Пушкина",      "21", null,  null,  esc, tc);
        Address a8 = addr("пр. Победы",       "44", null,  "9",   esc, school);
        Address a9 = addr("ул. Кирова",       "17", "2",   "33",  esc, energo);
        Address a10= addr("ул. Строителей",   "9",  null,  "11",  esc, hosp);

        // ── Addresses — others ──────────────────────────────────────────────────
        Address ag1 = addr("пр. Октября",      "33", null, "5",  gaz, gazt);
        Address ag2 = addr("ул. Молодёжная",   "8",  "1",  null, gaz, mkd);
        Address as1 = addr("ул. Новостроек",   "1",  "А",  "42", str, zhk);
        Address am1 = addr("пр. Автозаводцев", "100",null, null, mw,  zavod);
        Address am2 = addr("ул. Промышленная", "55", "3",  null, mw,  mega);
        Address an1 = addr("School lane",      "7",  null, null, ng,  nschool);

        // ── Tasks — esc-ural (completed) ────────────────────────────────────────
        task(TaskType.INSPECTION,  esc, a1,  TaskStatus.COMPLETED, LocalDate.now().minusDays(15), kozlov);
        task(TaskType.REPLACEMENT, esc, a3,  TaskStatus.COMPLETED, LocalDate.now().minusDays(10), fedorova);
        task(TaskType.INSPECTION,  esc, a4,  TaskStatus.COMPLETED, LocalDate.now().minusDays(8),  kozlov);
        task(TaskType.REPLACEMENT, esc, a2,  TaskStatus.COMPLETED, LocalDate.now().minusDays(5),  volkov);
        task(TaskType.INSPECTION,  esc, a9,  TaskStatus.COMPLETED, LocalDate.now().minusDays(3),  fedorova);

        // ── Tasks — esc-ural (in_progress) ──────────────────────────────────────
        task(TaskType.REPLACEMENT, esc, a5,  TaskStatus.IN_PROGRESS, LocalDate.now().plusDays(2),  kozlov);
        task(TaskType.INSPECTION,  esc, a6,  TaskStatus.IN_PROGRESS, LocalDate.now().plusDays(3),  fedorova);
        task(TaskType.REPLACEMENT, esc, a8,  TaskStatus.IN_PROGRESS, LocalDate.now().plusDays(4),  volkov);

        // ── Tasks — esc-ural (pending) ───────────────────────────────────────────
        task(TaskType.INSPECTION,  esc, a7,  TaskStatus.PENDING, LocalDate.now().plusDays(5),  kozlov);
        task(TaskType.REPLACEMENT, esc, a10, TaskStatus.PENDING, LocalDate.now().plusDays(7),  fedorova);
        task(TaskType.INSPECTION,  esc, a1,  TaskStatus.PENDING, LocalDate.now().plusDays(9),  volkov);
        task(TaskType.REPLACEMENT, esc, a4,  TaskStatus.PENDING, LocalDate.now().plusDays(12), kozlov);
        task(TaskType.REPLACEMENT, esc, a8,  TaskStatus.CANCELED, LocalDate.now().minusDays(1), volkov);

        // ── Tasks — other tenants ────────────────────────────────────────────────
        task(TaskType.INSPECTION,  gaz, ag1, TaskStatus.COMPLETED,   LocalDate.now().minusDays(7),  novikov);
        task(TaskType.REPLACEMENT, gaz, ag2, TaskStatus.IN_PROGRESS, LocalDate.now().plusDays(3),   novikov);
        task(TaskType.INSPECTION,  gaz, ag1, TaskStatus.PENDING,     LocalDate.now().plusDays(6),   novikov);
        task(TaskType.INSPECTION,  str, as1, TaskStatus.PENDING,     LocalDate.now().plusDays(4),   titov);
        task(TaskType.REPLACEMENT, str, as1, TaskStatus.COMPLETED,   LocalDate.now().minusDays(2),  titov);
        task(TaskType.REPLACEMENT, mw,  am1, TaskStatus.COMPLETED,   LocalDate.now().minusDays(20), gorbunov);
        task(TaskType.INSPECTION,  mw,  am2, TaskStatus.IN_PROGRESS, LocalDate.now().plusDays(1),   ershov);
        task(TaskType.REPLACEMENT, mw,  am1, TaskStatus.PENDING,     LocalDate.now().plusDays(8),   gorbunov);
        task(TaskType.INSPECTION,  ng,  an1, TaskStatus.PENDING,     LocalDate.now().plusDays(3),   null);

        log.info("=== Demo data created successfully ===");
        log.info("Android login: tenantCode=esc-ural  login=kozlov_d  password=demo123");
        log.info("Admin panel:   superadmin / admin123");
        log.info("Tenant admin:  tenantCode=esc-ural  login=ivanov_p  password=admin123");
    }

    // ── Helpers ────────────────────────────────────────────────────────────────

    private Tenant save(Tenant t) { return tenantRepo.save(t); }

    private static Tenant t(String name, String code, TenantPlan plan) {
        return Tenant.builder().name(name).code(code).plan(plan).active(true).build();
    }

    private static User user(String username, String password, String fullName,
                              UserRole role, Tenant tenant) {
        return User.builder()
                .username(username).password(password)
                .fullName(fullName).role(role).tenant(tenant)
                .build();
    }

    private static Organization org(String name, OrganizationType type, Tenant tenant) {
        return Organization.builder().name(name).type(type).tenant(tenant).build();
    }

    private Address addr(String street, String house, String building,
                         String apartment, Tenant tenant, Organization consumer) {
        return addressRepo.save(Address.builder()
                .street(street).house(house).building(building)
                .apartment(apartment).tenant(tenant).consumer(consumer)
                .build());
    }

    private void task(TaskType type, Tenant tenant, Address address,
                      TaskStatus status, LocalDate dueDate, User assignee) {
        taskRepo.save(Task.builder()
                .type(type).tenant(tenant).address(address)
                .status(status).dueDate(dueDate).assignee(assignee)
                .build());
    }
}
