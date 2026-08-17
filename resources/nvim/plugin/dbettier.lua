local package_name = "postgres-language-server"
local server_command = package_name
local lsp_enabled = false

local function notify(message, level)
  vim.schedule(function()
    vim.notify(message, level or vim.log.levels.INFO, { title = "dbettier" })
  end)
end

local function find_mason_server()
  local ok, settings = pcall(require, "mason.settings")
  if not ok then
    return nil
  end

  local path = settings.current.install_root_dir .. "/bin/" .. package_name
  if vim.fn.executable(path) == 1 then
    return path
  end

  return nil
end

local function start_lsp(bufnr)
  if not vim.api.nvim_buf_is_valid(bufnr) or vim.bo[bufnr].filetype ~= "sql" then
    return
  end

  local root_dir = vim.env.DBETTIER_PGLS_ROOT
  if not root_dir or root_dir == "" then
    local buffer_path = vim.api.nvim_buf_get_name(bufnr)
    root_dir = buffer_path ~= "" and vim.fs.dirname(buffer_path) or vim.fn.getcwd()
  end

  vim.lsp.start({
    name = "postgres_lsp",
    cmd = { server_command, "lsp-proxy" },
    root_dir = root_dir,
  }, {
    bufnr = bufnr,
  })
end

local function enable_lsp()
  if lsp_enabled then
    return
  end
  lsp_enabled = true

  local group = vim.api.nvim_create_augroup("dbettier_postgres_lsp", { clear = true })
  vim.api.nvim_create_autocmd("FileType", {
    group = group,
    pattern = "sql",
    callback = function(args)
      start_lsp(args.buf)
    end,
  })

  for _, bufnr in ipairs(vim.api.nvim_list_bufs()) do
    start_lsp(bufnr)
  end
end

local function use_installed_server()
  if vim.fn.executable(package_name) == 1 then
    server_command = package_name
    enable_lsp()
    return true
  end

  local mason_server = find_mason_server()
  if mason_server then
    server_command = mason_server
    enable_lsp()
    return true
  end

  return false
end

if use_installed_server() then
  return
end

local has_mason, registry = pcall(require, "mason-registry")
if not has_mason then
  notify(
    package_name .. " was not found and mason.nvim is unavailable. Install the server system-wide or with Mason.",
    vim.log.levels.WARN
  )
  return
end

notify("Installing " .. package_name .. " with Mason...")

registry.refresh(function()
  vim.schedule(function()
    if not registry.has_package(package_name) then
      notify("Mason does not provide the package " .. package_name .. ".", vim.log.levels.ERROR)
      return
    end

    local package = registry.get_package(package_name)
    if package:is_installed() then
      if not use_installed_server() then
        notify("Mason installed " .. package_name .. ", but its executable could not be found.", vim.log.levels.ERROR)
      end
      return
    end

    local ok, err = pcall(function()
      package:install({}, function(success)
        vim.schedule(function()
          if success and use_installed_server() then
            notify(package_name .. " installed successfully.")
          else
            notify("Failed to install " .. package_name .. ". Check :MasonLog for details.", vim.log.levels.ERROR)
          end
        end)
      end)
    end)

    if not ok then
      notify("Failed to start Mason installation: " .. tostring(err), vim.log.levels.ERROR)
    end
  end)
end)
