* Таск-трекер должен быть отдельным дашбордом и доступен всем сотрудникам компании UberPopug Inc.   
  **Actor: User (roles: any)  
  Query: Get user tasks**


* Авторизация в таск-трекере должна выполняться через общий сервис авторизации UberPopug Inc.  
  **Actor: User (roles: any)  
  Command: Authorize  
  Data: Biometry  
  Event: User.LoggedIn**


* Новые таски может создавать кто угодно (администратор, начальник, разработчик, менеджер и любая другая роль). У задачи должны быть описание, статус (выполнена или нет) и попуг, на которого заассайнена задача.  
  **Actor: User (roles: any)  
  Command: Create task  
  Data: Task  
  Event: Tasks.Created**


* Менеджеры или администраторы должны иметь кнопку «заассайнить задачи», которая возьмёт все открытые задачи и рандомно заассайнит каждую на любого из сотрудников (кроме менеджера и администратора).    
  **Actor: User (roles: admin,manager)   
  Command: Shuffle tasks  
  Data: ???  
  Event: Tasks.Shuffled**


* Каждый сотрудник должен иметь возможность видеть в отдельном месте список заассайненных на него задач.  
  **Actor: User (roles: any)  
  Query: Get user tasks**


* Каждый сотрудник должен иметь возможность отметить задачу выполненной.  
  **Actor: User (roles: any)  
  Command: Complete task  
  Data: Task ID  
  Event: Tasks.Completed**


* Авторизация в дешборде аккаунтинга должна выполняться через общий сервис аутентификации UberPopug Inc.  
  **Actor: User (roles: admin,accountant)  
  Command: Authorize  
  Data: Biometry  
  Event: User.LoggedIn**


* Цены на задачу определяется единоразово, в момент появления в системе (можно с минимальной задержкой)  
  **Actor: "Tasks.Created" event  
  Command: Calculate task value  
  Data: Task  
  Event: Tasks.ValueUpdated**


* Деньги списываются сразу после ассайна на сотрудника.  
  **Actor: "Tasks.ValueUpdated", "Tasks.Shuffled" events  
  Command: Withdraw money from account  
  Data: Task, User  
  Event: Accounts.MoneyWithdrawn**


* Деньги начисляются после выполнения задачи.  
  **Actor: "Tasks.Completed" event  
  Command: Add money to account  
  Data: Task, User  
  Event: Accounts.MoneyAdded**