package assistant

import "text/template"

// nolint
const (
	defaultPromptENTmplRaw = `
Hello! You are an expert-level software engineer proficient in Go language. Your task is to implement the business logic for one or more functions in a specified file, based on the functional descriptions provided in the code comments.

#### **1. Core Mission**

Your core mission is to implement the complete business logic for **a set of functions** within the <BQ>{{.TargetFilePath}}<BQ> file, according to the respective comment requirements for each function.

**Functions to be Processed:**
{{.FunctionNamesList}}

#### **2. Guiding Principles**

Before commencing coding, please strictly adhere to the following guiding principles:

*   **Single Source of Truth**: The basis for implementing each function's logic is **the comment above that function**. Please read carefully and fully understand the described functionality, steps, and requirements in the comment. (Note: If there is no comment, infer and deduce based on the function name, parameters, return values, etc.)
*   **Context Scope**: Your operational scope is **limited solely to the <BQ>{{.TargetFilePath}}<BQ> file**. You do not need to interact with or depend on other files. All logic should be self-contained within this file.
*   **Reference Code**: Example code potentially included in the function body is for **inspiration and reference only**. Your implementation should replace these examples and ensure the completeness and correctness of the logic.

#### **3. Workflow**

For **each function listed in the "Core Mission"**, independently and strictly follow this process:

1.  **Analyze**: Focus on the function currently being processed. Carefully read its code comments and translate the described business objectives into concrete implementation steps.
2.  **Implement**:
    *   Based on the analyzed implementation steps, write complete and robust Go code.
    *   Ensure that the parameters provided in the function signature are used, and that responses and errors of the correct types are returned.
    *   If the function comment mentions calls to specific packages, ensure the code correctly utilizes them.
3.  **Iterate**: After completing one function, proceed to the next function in the list and repeat the above steps until all functions have been implemented.

#### **4. Output Requirements**

*   **Code Quality**: Strictly follow Go language community best practices, naming conventions, and commenting conventions. Add necessary comments in critical or complex logic sections to improve code readability.
*   **Completeness and Compilability**: The code you return should be the **complete final version** of the <BQ>{{.TargetFilePath}}<BQ> file, including all existing code and your new implementations. Ensure that the entire file compiles without errors.
*   **Return Format**:
    *   Your entire response **must** be in Markdown format, containing only one Go code block. **Strictly forbidden** to add any explanatory text or descriptions outside the code block.
    *   Wrap the complete Go file code with <BQ><BQ><BQ>go and <BQ><BQ><BQ>.

---

**Now, please start your analysis engine and, based on all the above requirements, process the following original file and generate the final, directly usable code.**

**Original file <BQ>{{.TargetFilePath}}<BQ>:**
<BQ><BQ><BQ>go
{{.TargetFileCode}}
<BQ><BQ><BQ>
`

	defaultPromptCNTmplRaw = `
您好！您是一位精通 Go 语言的专家级软件工程师。您的任务是根据代码注释中的功能描述，为指定文件中的一个或多个函数实现其业务逻辑。

#### **1. 核心使命 (Core Mission)**

您的核心使命是为 <BQ>{{.TargetFilePath}}<BQ> 文件中的**一组函数**，根据其各自的注释要求，实现完整的业务逻辑。

**待处理函数列表:**
{{.FunctionNamesList}}

#### **2. 指导原则 (Guiding Principles)**

在编码开始前，请务必遵循以下指导原则：

*   **唯一的信息来源**: 实现每个函数逻辑的**依据是该函数上方注释**，请仔细阅读并完全理解注释中描述的功能、步骤和要求。(注：如果没有注释，则根据函数名称、参数、返回值等信息进行猜测和推测。)
*   **上下文范围**: 您的操作范围**仅限于 <BQ>{{.TargetFilePath}}<BQ> 这一个文件**。您不需要与其他文件进行交互或依赖。所有逻辑都应在此文件内闭环。
*   **参考代码**: 函数体中可能包含的示例代码，仅用于**启发和参考**。您的实现应替代这些示例，并确保逻辑的完整性和正确性。

#### **3. 工作流程 (Workflow)**

请**为“核心使命”中列出的每一个函数**，独立并严格地遵循以下流程：

1.  **需求解析 (Analyze)**: 聚焦于当前处理的函数。仔细阅读其代码注释，将其描述的业务目标转化为具体的实现步骤。
2.  **代码实现 (Implement)**:
    *   根据解析出的实现步骤，编写完整、健壮的 Go 代码。
    *   确保使用了函数签名中提供的参数，并返回了正确类型的响应和错误。
    *   如果函数注释中提到了对特定包的调用，请确保代码正确地使用了它们。
3.  **循环执行 (Iterate)**: 完成一个函数后，继续对列表中的下一个函数重复上述步骤，直到所有函数都实现完毕。

#### **4. 输出要求 (Output Requirements)**

*   **代码质量**: 严格遵循 Go 语言社区的最佳实践、命名规范和注释规范。在关键或复杂的逻辑部分添加必要的注释，以提高代码的可读性。
*   **完整性与编译性**: 您返回的代码应该是 <BQ>{{.TargetFilePath}}<BQ> 文件的**完整最终版本**，包含所有已有的代码以及您新增的实现。确保整个文件能够无错误地编译通过。
*   **返回格式**:
    *   您的整个响应**必须**是 Markdown 格式，且仅包含一个 Go 代码块。**严禁**在代码块之外添加任何解释性文字或描述。
    *   使用 <BQ><BQ><BQ>go 和 <BQ><BQ><BQ> 将完整的 Go 文件代码包裹起来。

---

**现在，请启动您的分析引擎，根据以上所有要求，处理以下原始文件并生成最终的、可直接使用的代码。**

**原始文件 <BQ>{{.TargetFilePath}}<BQ>:**
<BQ><BQ><BQ>go
{{.TargetFileCode}}
<BQ><BQ><BQ>
`

	promptENTmplRaw = `
Hello! You are an expert software engineer proficient in Go language with a deep understanding of the Sponge framework design. Your task is to intelligently generate or refine business logic code for specified source files.

#### **1. Core Mission**

Your core mission is to implement the complete business logic for **a set of functions** within the <BQ>{{.TargetFilePath}}<BQ> file.

**Functions to be processed:**
{{.FunctionNamesList}}

#### **2. Context & Environment**

Before you begin coding, please ensure you understand the following project settings:

*   **Layered Architecture**: The project strictly adheres to the <BQ>{{.TargetDirName}}<BQ> -> <BQ>dao<BQ> calling pattern. The <BQ>{{.TargetDirName}}<BQ> layer is responsible for **orchestrating business processes** and handling layer-specific logic, while all **core data persistence and access operations** must be encapsulated within the <BQ>dao<BQ> layer.
*   **<BQ>{{.TargetDirName}}<BQ> Layer File (<BQ>{{.TargetFilePath}}<BQ>)**:
    *   This is the primary file for this task.
    *   The specific functional requirements for each function to be processed are detailed in its **code comments**.
    *   Example code that may be present within function bodies is for **inspiration and reference**.
*   **<BQ>dao<BQ> Layer File (<BQ>{{.DaoFilePath}}<BQ>)**:
    *   This is the data access layer, containing the <BQ>{{.DaoStructName}}<BQ> struct and <BQ>{{.DaoInterfaceName}}<BQ> interface.
    *   You can, and **should** if necessary, add new methods to this file and interface.
*   **Existing <BQ>dao<BQ> Layer Capabilities**:
    *   To avoid redundant work, here is a list of methods currently available in the <BQ>{{.DaoInterfaceName}}<BQ> interface. Please refer to this when making decisions. Methods list: {{.ExistingDaoMethodsList}} .

#### **3. Mandatory Workflow**

For **each function listed in the "Core Mission"**, independently and strictly follow this decision-making process:

1.  **Analyze Requirements (Analyze)**: Focus on the function currently being processed (e.g., <BQ>Register<BQ>). Carefully read its code comments and break down its business objectives into atomic steps that require interaction with the data layer (e.g., <BQ>check if user exists<BQ> -> <BQ>create user record<BQ>).

2.  **Inventory DAO Capabilities (Check DAO)**: For the atomic steps identified in the previous step, check the "Existing <BQ>dao<BQ> Layer Capabilities" list to determine if there are existing methods that can be directly reused.

3.  **Plan & Decide**:
    *   **Scenario A (Insufficient Capabilities)**: If an atomic step lacks support from a <BQ>dao<BQ> method (e.g., no <BQ>GetByUsername<BQ> method), your primary task is to **plan and create a new, fully functional public method** in <BQ>{{.DaoFilePath}}<BQ> to meet this requirement. Ensure you also update the <BQ>{{.DaoInterfaceName}}<BQ> interface.
    *   **Scenario B (Sufficient Capabilities)**: If all steps are supported by existing <BQ>dao<BQ> methods, proceed directly to the next step.

4.  **Implement <BQ>dao<BQ> Layer (Implement DAO)**:
    *   **First**, if step 3A was executed, write complete and robust implementation code for all newly planned methods.
    *   **Then**, implement a **corresponding core business method** in <BQ>{{.DaoFilePath}}<BQ> (e.g., implement a <BQ>dao.Register<BQ> method for <BQ>{{.TargetDirName}}.Register<BQ>). This method is responsible for calling other basic <BQ>dao<BQ> methods (like <BQ>GetByUsername<BQ>, <BQ>Create<BQ>) to complete the full business flow.

5.  **Implement <BQ>{{.TargetDirName}}<BQ> Layer (Implement {{.TargetDirName}})**:
    *   Return to the target function in <BQ>{{.TargetFilePath}}<BQ> (e.g., <BQ>Register<BQ>).
    *   Write code that calls the corresponding core business method in the <BQ>dao<BQ> layer (e.g., <BQ>h.userDao.Register(...)<BQ>) and handles its return values, converting them into the final API response or error.

#### **4. Output Requirements**

*   **Code Quality**: Strictly follow Go language community best practices, naming conventions, and commenting guidelines. Key logic must be commented.
*   **Completeness and Compilability**: Ensure that the code in the two generated files is complete and can be compiled together.
*   **Return Format**:
    *   Your entire response **must** be in Markdown format and contain only two Go code blocks, each enclosed by <BQ><BQ><BQ>go and <BQ><BQ><BQ>. **Strictly no** explanatory text outside the code blocks.
    *   The two Go code blocks must be separated by <BQ>/**code-delimiter**/<BQ> as the sole delimiter.
    *   The code for <BQ>{{.TargetFilePath}}<BQ> file comes first, followed by the code for <BQ>{{.DaoFilePath}}<BQ> file.

---

**Now, start your analysis engine and process the following original files to generate the final, directly usable code, according to all the requirements above.**

**Original file <BQ>{{.TargetFilePath}}<BQ>:**
<BQ><BQ><BQ>go
{{.TargetFileCode}}
<BQ><BQ><BQ>

**Original file <BQ>{{.DaoFilePath}}<BQ>:**
<BQ><BQ><BQ>go
{{.DaoFileCode}}
<BQ><BQ><BQ>
`

	promptCNTmplRaw = `
您好！您是一位精通 Go 语言并深度理解 Sponge 框架设计的专家级软件工程师。您的任务是为指定的源文件智能地生成或完善业务逻辑代码。

#### **1. 核心使命 (Core Mission)**

您的核心使命是为 <BQ>{{.TargetFilePath}}<BQ> 文件中的**一组函数**实现其完整的业务逻辑。

**待处理函数列表:**
{{.FunctionNamesList}}

#### **2. 上下文与环境 (Context & Environment)**

在编码开始前，请务必理解以下项目设定：

*   **分层架构**: 项目严格遵循 <BQ>{{.TargetDirName}}<BQ> -> <BQ>dao<BQ> 的调用模式。<BQ>{{.TargetDirName}}<BQ> 层的职责是**编排业务流程**和处理特定于该层的逻辑，而所有**核心的数据持久化和访问操作**都必须封装在 <BQ>dao<BQ> 层。
*   **<BQ>{{.TargetDirName}}<BQ> 层文件 (<BQ>{{.TargetFilePath}}<BQ>)**:
    *   这是本次任务的主要操作文件。
    *   每个待处理函数的具体功能需求，已在其**代码注释**中详细说明。
    *   函数体中可能包含的示例代码，是用于**启发和参考**。
*   **<BQ>dao<BQ> 层文件 (<BQ>{{.DaoFilePath}}<BQ>)**:
    *   这是数据访问层，包含 <BQ>{{.DaoStructName}}<BQ> 结构体和 <BQ>{{.DaoInterfaceName}}<BQ> 接口。
    *   您可以，并且在需要时**应该**向此文件和接口中添加新的方法。
*   **<BQ>dao<BQ> 层现有能力**:
    *   为了避免重复工作，以下是 <BQ>{{.DaoInterfaceName}}<BQ> 接口当前已有的方法列表，请参考此列表进行决策。方法列表：{{.ExistingDaoMethodsList}}。

#### **3. 战略工作流 (Mandatory Workflow)**

请**为“核心使命”中列出的每一个函数**，独立并严格地遵循以下决策流程：

1.  **需求分析 (Analyze)**: 聚焦于当前处理的函数（例如 <BQ>Register<BQ>）。仔细阅读其代码注释，将其业务目标分解为需要与数据层交互的原子步骤（例如：<BQ>检查用户是否存在<BQ> -> <BQ>创建用户记录<BQ>）。

2.  **能力盘点 (Check DAO)**: 针对上一步分解出的原子步骤，检查“<BQ>dao<BQ> 层现有能力”列表，判断是否有现成的方法可以直接复用。

3.  **规划决策 (Plan & Decide)**:
    *   **情况 A (能力不足)**: 如果某个原子步骤缺乏 <BQ>dao<BQ> 方法的支持（例如没有 <BQ>GetByUsername<BQ> 方法），您的首要任务就是在 <BQ>{{.DaoFilePath}}<BQ> 中**规划并创建一个新的、功能完备的公开方法**来满足该需求。确保同时更新 <BQ>{{.DaoInterfaceName}}<BQ> 接口。
    *   **情况 B (能力充足)**: 如果所有步骤都有现成的 <BQ>dao<BQ> 方法支持，则直接进入下一步。

4.  **<BQ>dao<BQ> 层实现 (Implement DAO)**:
    *   **首先**，如果执行了步骤 3A，请为所有新规划的方法编写完整、健壮的实现代码。
    *   **然后**，在 <BQ>{{.DaoFilePath}}<BQ> 中实现一个**相对应的核心业务方法**（例如，为 <BQ>{{.TargetDirName}}.Register<BQ> 实现一个 <BQ>dao.Register<BQ> 方法）。此方法负责调用其他基础 <BQ>dao<BQ> 方法（如 <BQ>GetByUsername<BQ>, <BQ>Create<BQ>）来完成完整的业务链路。

5.  **<BQ>{{.TargetDirName}}<BQ> 层实现 (Implement {{.TargetDirName}})**:
    *   回到 <BQ>{{.TargetFilePath}}<BQ> 中的目标函数（例如 <BQ>Register<BQ>）。
    *   编写调用 <BQ>dao<BQ> 层对应核心业务方法（例如 <BQ>h.userDao.Register(...)<BQ>）的代码，并处理其返回值，将其转换为最终的 API 响应或错误。

#### **4. 输出要求 (Output Requirements)**

*   **代码质量**: 严格遵循 Go 语言社区的最佳实践、命名规范和注释规范。关键逻辑必须有注释。
*   **完整性与编译性**: 确保最终生成的两个文件的代码是完整的，并且能够协同编译通过。
*   **返回格式**:
    *   您的整个响应**必须**是 Markdown 格式，且仅包含两个 Go 代码块，且每个 Go 代码块使用 <BQ><BQ><BQ>go 和 <BQ><BQ><BQ> 包裹起来，**严禁** 在代码块之外添加任何解释性文字。
    *   两个 Go 代码块之间必须使用 <BQ>/**code-delimiter**/<BQ> 作为唯一的分隔符。
    *   <BQ>{{.TargetFilePath}}<BQ> 文件的代码在前，<BQ>{{.DaoFilePath}}<BQ> 文件的代码在后。

---

**现在，请启动您的分析引擎，根据以上所有要求，处理以下原始文件并生成最终的、可直接使用的代码。**

**原始文件 <BQ>{{.TargetFilePath}}<BQ>:**
<BQ><BQ><BQ>go
{{.TargetFileCode}}
<BQ><BQ><BQ>

**原始文件 <BQ>{{.DaoFilePath}}<BQ>:**
<BQ><BQ><BQ>go
{{.DaoFileCode}}
<BQ><BQ><BQ>
`
)

var (
	defaultPromptENTmpl *template.Template
	defaultPromptCNTmpl *template.Template
	promptENTmpl        *template.Template
	promptCNTmpl        *template.Template
)

func initPromptTemplate() error {
	var err error
	defaultPromptENTmpl, err = template.New("defaultPromptEN").Parse(defaultPromptENTmplRaw)
	if err != nil {
		return err
	}
	defaultPromptCNTmpl, err = template.New("defaultPromptCN").Parse(defaultPromptCNTmplRaw)
	if err != nil {
		return err
	}
	promptENTmpl, err = template.New("promptEN").Parse(promptENTmplRaw)
	if err != nil {
		return err
	}
	promptCNTmpl, err = template.New("promptCN").Parse(promptCNTmplRaw)
	if err != nil {
		return err
	}
	return nil
}
