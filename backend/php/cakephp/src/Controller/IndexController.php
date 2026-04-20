<?php
declare(strict_types=1);

namespace App\Controller;

class IndexController extends AppController
{
    public function browser()
    {
        $this->responda(200, (object)[], 'Aplicação deverá ser acessada via /api');
    }

    public function index()
    {
        $this->responda(200, (object)[], 'Bem vindo à API da Igreja Aberta');
    }
}